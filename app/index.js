const {app, BrowserWindow} = require('electron');

function createWindow() {
  let win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      // TODO: How to handle data? Web worker? preload + disable node integration?
      nodeIntegration: false,
    }
  });

  win.once('ready-to-show', () => {
    win.show();
  });

  // TODO: How to handle generating and sharing a certificate so this can use SSL.
  win.loadURL("http://localhost:10808/index.html", {
    show: false,
  });

  win.webContents.openDevTools();

  return win;
}

async function startServer(username, password) {
  const defaults = {
    cwd: undefined,
    env: process.env,
  };

  // TODO: Receive these from the command-line.
  // TODO: Use go build instead of go run.
  let args = [
    'run',
    './cmd/fortnight/main.go',
    '-log-level=debug',
    '-insecure-skip-verify',
    '-server=localhost:18443',
  ];

  if (username !== undefined && password !== undefined) {
    args.push(`-credentials=${username}:${password}`)
  }

  // TODO: These promises are a mess, but it seems to work.
  // The intent is to create a promise that resolves after the process
  // is started, returning another promise that resolves after the process
  // exits.
  return new Promise((resolve, reject) => {
    const {spawn} = require('child_process');
    const server = spawn('vgo', args, defaults);

    let started = false;
    server.stderr.on('data', (data) => {
      if (data.indexOf('listening on') >= 0) {
        started = true;
      }
    });

    let stopped = false;
    let failed = false;
    server.on('close', (code) => {
      console.log(`child process exited with code ${code}`);
      stopped = true;
      if (code > 0) {
        failed = true;
      }
    });

    const finish = new Promise((resolve, reject) => {
      let iv;
      // Check if the process is stopped at regular intervals.
      iv = setInterval(function() {
        if (stopped === true) {
          if (failed === true) {
            console.log('server process stopped with failure');
            reject();
          } else {
            console.log('server process stopped cleanly');
            resolve();
          }
          clearInterval(iv);
        }
      }, 5000);
    });

    let iv;
    // Check if the process is started at regular, very short, intervals.
    iv = setInterval(function() {
      if (started === true) {
        resolve();
        clearInterval(iv);
      }
    }, 100);

    return finish;
  });
}

// TODO: Handle this in the app?
const username = process.argv[2] || undefined;
const password = process.argv[3] || undefined;

app.on('ready', async () => {
  return startServer(username, password)
      .then(() => createWindow())
      .catch(reason => console.log(reason));
});
