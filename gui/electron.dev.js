const { app, BrowserWindow, shell, session } = require('electron');

let win

function createWindow() {
    let height = 500;

    win = new BrowserWindow({
        width: 900,
        height,
        autoHideMenuBar: true,
        resizable: false
    });

    win.setSize(900, 760);
    win.webContents.openDevTools();

    win.loadURL("http://localhost:3002/")
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
    })
}

app.on('ready', () => {
    setupCookie("enginePort", "3000");
    setupCookie("engineAddr", "localhost");
    setupCookie("systemLocale", app.getLocale());
    createWindow();
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

function setupCookie(name, value) {
    console.log("Registering Cookie: ", name, value);
    const cookie = { url: 'http://localhost:3002', name, value };
    session.defaultSession.cookies.set(cookie, (error) => {
        if (error) console.error(error);
    });
}