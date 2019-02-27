const { app, BrowserWindow, protocol } = require('electron');
const express = require('express');
const path = require('path');
const http = require('http');

let win

function createWindow() {
    let height = 500;

    win = new BrowserWindow({
        width: 900,
        height,
        autoHideMenuBar: true,
        resizable: false
    })

    win.setSize(900, 760)
    win.webContents.openDevTools();

    win.loadURL("http://localhost:3002/index.html")
    win.focus();

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