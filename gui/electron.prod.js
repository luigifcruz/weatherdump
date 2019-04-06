const { app, BrowserWindow, protocol, shell, dialog, session } = require('electron');
const { spawn } = require('child_process');
const getPort = require('get-port');
const express = require('express');
const http = require('http');
const path = require('path');

const serve = new express();

let win, cli, server, electronPort, enginePort = null;

function createWindow() {
    let height = 500;

    if (process.platform == 'darwin' || process.platform == 'win32') {
        height = 525;
    }

    win = new BrowserWindow({
        width: 900,
        height,
        autoHideMenuBar: true,
        resizable: false,
        backgroundColor: "#212330"
    });

    if (process.env.NODE_ENV == 'debug') {
        win.setSize(900, 760);
        win.webContents.openDevTools();
    }

    win.loadURL("http://localhost:"+electronPort)
    win.focus();

    win.webContents.on('new-window', function(e, payload) {
        e.preventDefault();
        const url = new URL(payload);
        if (url.hostname == "localhost") {
            shell.openItem(url.pathname);
        } else {
            shell.openExternal(url.href);
        }
    });

    win.on('closed', () => {
        win = null
    });
}

app.on('will-quit', () => {
    console.log("Safely quiting...")
    if (cli) {
        cli.stdin.pause();
        cli.kill();
        cli = null;
    }
    server.close();
});

app.on('ready', () => {
    (async () => {
        enginePort = await getPort({port: getPort.makeRange(3050, 3150)});
        electronPort = await getPort({port: getPort.makeRange(3100, 3150)});

        setupCookie("engineAddr", "localhost");
        setupCookie("enginePort", enginePort.toString());
        setupCookie("electronPort", electronPort.toString());
        setupCookie("systemLocale", app.getLocale());

        await startServer();
        startEngine();
        
        createWindow();
    })();

    if (process.platform === 'win32') {
        app.setAppUserModelId('com.osp.weatherdump');
    }
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

app.on('activate', () => {
    if (win === null) {
        createWindow();
    }
});

function getBinaryPath() {
    let binaryName = "weatherdump"

    if (process.platform == "win32") {
        binaryName.concat(".exe")
    }

    return path.join(__dirname, "..", "app", "engine", binaryName)
}

function startServer() {
    return new Promise((resolve, reject) => {
        serve.use('/', express.static(path.join(__dirname, "..", "app", "gui")))
        server = http.createServer(serve).listen({
            host: 'localhost',
            port: electronPort,
            exclusive: true
        }, () => {
            console.log("Server started listening to port %d.", electronPort);
            resolve();
	    });
    });
}

function startEngine() {
    cli = spawn(getBinaryPath(), ['remote', enginePort.toString(), electronPort.toString()]);

    cli.on('exit', (code) => {
        if (cli != null) {
            dialog.showErrorBox(
                "Unexpected Engine Crash",
                "Something has gone terribly wrong with the WeatherDump engine. Please, report this error to @luigifcruz at Twitter."
            );
            cli = null;
            app.quit();
        }
    });
}

function setupCookie(name, value) {
    const cookie = { url: "http://localhost:" + electronPort, name, value };

    console.log("Registering Cookie: ", name, value);
    session.defaultSession.cookies.set(cookie, (error) => {
        if (error) console.error(error);
    });
}
