const { app, BrowserWindow } = require('electron')

let win

function createWindow() {
    let height = 500;

    if (process.platform == 'darwin') {
        height = 520
    }

    win = new BrowserWindow({
        width: 900,
        height,
        autoHideMenuBar: true,
        resizable: false
    })

    if (process.env.NODE_ENV === "development") {
        win.setSize(900, 750)
        win.webContents.openDevTools();
    }

    win.loadURL('http://localhost:3002');

    win.on('closed', () => {
        win = null
    })
}

app.on('ready', createWindow)

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit()
    }
})

app.on('activate', () => {
    if (win === null) {
        createWindow()
    }
})