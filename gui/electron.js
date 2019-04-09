const { app, BrowserWindow, protocol, shell } = require('electron');
const path = require('path');
const url = require('url');

let win = null;

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

    if (process.env.NODE_ENV === 'development') {
        console.log("Development mode enabled...");
        win.setSize(900, 760);
        win.webContents.openDevTools();
        win.loadURL("http://localhost:3002/index.html");
    } else {
        console.log("Production mode enabled...");
        win.loadURL(url.format({
            pathname: 'index.html',
            protocol: 'file',
            slashes: true
        }));
    }

    win.focus();
    win.on('closed', () => {
        win = null;
    });

    win.webContents.on('new-window', function(e, payload) {
        e.preventDefault();
        const url = new URL(payload);
        if (url.hostname == "localhost") {
            shell.openItem(url.pathname);
        } else {
            shell.openExternal(url.href);
        }
    });
}

app.on('ready', () => {
    protocol.interceptFileProtocol('file', (request, callback) => {
        const url = request.url.substr(7)
        callback({ path: path.join(__dirname, "..", "app", "gui", url) })
    }, (err) => {
        if (err) console.error('Failed to register protocol')
    })

    createWindow();
    
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
