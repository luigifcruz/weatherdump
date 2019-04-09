import { spawn } from 'child_process';
import { remote } from 'electron';
import path from 'path';

class WeatherServer {
    constructor(port) {
        this.dirname = remote.app.getAppPath();
        this.port = port
        this.cli = null;
    }

    get isRunning() {
        return this.cli !== null;
    }

	getBinaryPath() {
		let binaryName = "weatherdump";
	
		if (process.platform == "win32") {
			binaryName.concat(".exe");
        }

        return path.join(this.dirname, "..", "app", "engine", binaryName);
	}
	
	startEngine() {
        this.cli = spawn(this.getBinaryPath(), ['remote', this.port]);
        
        this.cli.stdout.on('data', (data) => {
            console.log(data.toString());
        });

        this.cli.stderr.on('data', (data) => {
            console.error(data.toString());
        });

        this.cli.on('exit', (code) => {
            if (this.cli != null) {
                this.cli = null;
                this.reportCrash("Engine has exited with code " + code);
            }
        });
    }
    
    stopEngine() {
        if (!this.isRunning) {
            console.warn("Engine isn't running.");
            return;
        }

        this.cli.stdin.pause();
        this.cli.kill();
        this.cli = null;
    }
	
	reportCrash(crash) {
		console.error(crash);
	}
};

export default WeatherServer;