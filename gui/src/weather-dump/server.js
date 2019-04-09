import { spawn } from 'child_process';
import path from 'path';

class WeatherServer {
    constructor(port) {;
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
        
		return path.join(__dirname, "..", "app", "engine", binaryName);
	}
	
	startEngine() {
        this.cli = spawn(this.getBinaryPath(), ['remote', this.port]);
        
        this.cli.stdout.on('data', (data) => {
            process.stdout.write(data.toString());
        });

        this.cli.on('exit', () => {
            if (this.cli != null) {
                this.cli = null;
                this.reportCrash("Engine has crashed.");
            }
        });
    }
    
    stopEngine() {
        if (!this.isRunning) {
            console.error("Engine isn't running.");
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