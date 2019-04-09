const { app, BrowserWindow, shell } = require('electron');

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
        win = null;
    });
}

app.on('ready', () => {
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