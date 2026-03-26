const { app, BrowserWindow, Tray, Menu, ipcMain } = require('electron')
const path = require('path')
const Store = require('electron-store')

const store = new Store()

let mainWindow = null
let tray = null

// 窗口状态持久化键
const WINDOW_BOUNDS_KEY = 'window-bounds'

function createWindow() {
  // 读取保存的窗口位置
  const savedBounds = store.get(WINDOW_BOUNDS_KEY)

  const windowOptions = {
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    title: 'CloudIM',
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    }
  }

  // 恢复保存的窗口位置
  if (savedBounds) {
    windowOptions.x = savedBounds.x
    windowOptions.y = savedBounds.y
    windowOptions.width = savedBounds.width
    windowOptions.height = savedBounds.height
  }

  mainWindow = new BrowserWindow(windowOptions)

  // 开发环境加载 Vite 开发服务器
  const isDev = process.env.NODE_ENV !== 'production'
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    // 生产环境加载构建后的文件
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
  }

  // 窗口关闭时保存位置
  mainWindow.on('close', () => {
    if (!mainWindow.isDestroyed()) {
      const bounds = mainWindow.getBounds()
      store.set(WINDOW_BOUNDS_KEY, bounds)
    }
  })

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

function createTray() {
  // 使用默认图标（生产环境应该使用自定义图标）
  tray = new Tray(path.join(__dirname, 'icon.png'))

  const contextMenu = Menu.buildFromTemplate([
    {
      label: '显示主窗口',
      click: () => {
        if (mainWindow) {
          mainWindow.show()
          mainWindow.focus()
        }
      }
    },
    {
      label: '退出',
      click: () => {
        app.quit()
      }
    }
  ])

  tray.setToolTip('CloudIM')
  tray.setContextMenu(contextMenu)

  tray.on('click', () => {
    if (mainWindow) {
      if (mainWindow.isVisible()) {
        mainWindow.focus()
      } else {
        mainWindow.show()
      }
    }
  })
}

// 应用准备好时创建窗口
app.whenReady().then(() => {
  createWindow()
  createTray()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

// 所有窗口关闭时退出应用（macOS 除外）
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// IPC 处理
ipcMain.on('close-window', () => {
  if (mainWindow) {
    mainWindow.close()
  }
})

ipcMain.on('minimize-window', () => {
  if (mainWindow) {
    mainWindow.minimize()
  }
})

ipcMain.on('maximize-window', () => {
  if (mainWindow) {
    if (mainWindow.isMaximized()) {
      mainWindow.unmaximize()
    } else {
      mainWindow.maximize()
    }
  }
})
