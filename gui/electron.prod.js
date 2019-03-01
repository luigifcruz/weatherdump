const { app, BrowserWindow, protocol } = require('electron');
const express = require('express');
const path = require('path');
const http = require('http');
const { spawn } = require('child_process');
const url = require('url')

let win, cli

function createWindow() {
    let height = 500;

    if (process.platform == 'darwin' || process.platform == 'win32') {
        height = 525
    }

    win = new BrowserWindow({
        width: 900,
        height,
        autoHideMenuBar: true,
        resizable: false
    })

    if (process.env.NODE_ENV == 'debug') {
        win.setSize(900, 760)
        win.webContents.openDevTools();
    }

    win.loadURL(url.format({
        pathname: 'index.html',
        protocol: 'file',
        slashes: true
    }))

    win.focus();

    win.on('closed', () => {
        win = null
        cli.stdin.pause();
        cli.kill();
    })
}

app.on('ready', () => {
    protocol.interceptFileProtocol('file', (request, callback) => {
        const url = request.url.substr(7)
        callback({ path: path.join(__dirname, "..", "app", "gui", url) })
    }, (err) => {
        if (err) console.error('Failed to register protocol')
    })

    startEngine()
    createWindow()
})

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

function getBinaryPath() {
    let binaryName = "weatherdump"

    if (process.platform == "win32") {
        binaryName.concat(".exe")
    }

    return path.join(__dirname, "..", "app", "engine", binaryName)
}

function startEngine() {
    cli = spawn(getBinaryPath(), ['remote']);
}