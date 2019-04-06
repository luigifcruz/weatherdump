const { app, BrowserWindow, protocol, shell, dialog } = require('electron');
const path = require('path');
const getPort = require('get-port');
const { spawn } = require('child_process');
const url = require('url');

let win, cli = null;

function createWindow() {
    let height = 500;

    if (process.platform == 'darwin' || process.platform == 'win32') {
        height = 525;
    }

    win = new BrowserWindow({
        width: 900,
        height,
        autoHideMenuBar: true,
        resizable: false
    });

    if (process.env.NODE_ENV == 'debug') {
        win.setSize(900, 760);
        win.webContents.openDevTools();
    }

    win.loadURL(url.format({
        pathname: 'index.html',
        protocol: 'file',
        slashes: true
    }));

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
        if (cli) {
            cli.stdin.pause();
            cli.kill();
        }
    });
}

app.on('ready', () => {
    protocol.interceptFileProtocol('file', (request, callback) => {
        const url = request.url.substr(7);
        callback({ path: path.join(__dirname, "..", "app", "gui", url) });
    }, (err) => {
        if (err) console.error('Failed to register protocol');
    });

    (async () => {
        let enginePort = await getPort({port: getPort.makeRange(3050, 3150)});
        setupCookie("engineAddr", "localhost");
        setupCookie("enginePort", enginePort.toString());
        setupCookie("systemLocale", app.getLocale());
        startEngine(enginePort.toString());
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

function startEngine(enginePort) {
    cli = spawn(getBinaryPath(), ['remote', enginePort]);
    //cli = spawn("../dist/weatherdump-cli-linux-x64/weatherdump", ['remote']);

    cli.on('exit', (code) => { 
        cli = null;
        dialog.showErrorBox(
            "Unexpected Engine Crash",
            "Something has gone terribly wrong with the WeatherDump engine. Please, report this error to @luigifcruz, @lucasteske or @OpenSatProject at Twitter.");
        app.quit()
    });
}

function setupCookie(name, value) {
    console.log("Registering Cookie: ", name, value);
    const cookie = { url: 'http://localhost:3002', name, value };
    session.defaultSession.cookies.set(cookie, (error) => {
        if (error) console.error(error);
    });
}