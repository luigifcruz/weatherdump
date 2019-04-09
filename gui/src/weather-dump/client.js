import request from 'superagent';

class WeatherClient {
	constructor(engineAddr, enginePort) {
		this.serverAddress = `${engineAddr}:${enginePort}`;
	}

	get enginePort() {
		return this.enginePort;
	}

	get engineAddr() {
		return this.engineAddr;
	}

	get enginePath() {
		return this.serverAddress;
	}

	// Abort Task
	// Endpoint: /abort/:uuid

	abortTask(uuid) {
		return new Promise((resolve, reject) => {
			request
				.post(`http://${this.serverAddress}/abort/${uuid}`)
				.then(resolve)
				.catch(this.errorParser);
		});
	}

	// Start Processor
	// Endpoint: /start/processor
	// InputFile  string `schema:"inputPath,required"`
	// Datalink   string `schema:"datalink,required"`
	// Pipeline   string `schema:"pipeline,required"`
	// Manifest   string `schema:"manifest,required"`
	// OutputPath string `schema:"outputPath"`

	startProcessor(req) {
		return new Promise((resolve, reject) => {
			request
				.post(`http://${this.serverAddress}/start/processor`)
				.type('form')
				.send(req)
				.then((res) => {
					resolve({
						uuid: res.body.Code,
						outputPath: JSON.parse(res.body.Data).OutputPath
					});
				})
				.catch(this.errorParser)
		});
	}

	// Start Decoder
	// Endpoint: /start/decoder
	// InputFile  string `schema:"inputFile,required"`
	// Datalink   string `schema:"datalink,required"`
	// Decoder    string `schema:"decoder,required"`
	// OutputPath string `schema:"outputPath"`
	
	startDecoder(req) {
		return new Promise((resolve, reject) => {
			request
				.post(`http://${this.serverAddress}/start/decoder`)
				.type('form')
				.send(req)
				.then((res) => {
					resolve({
						uuid: res.body.Code,
						outputPath: JSON.parse(res.body.Data).OutputPath
					});
				})
				.catch(this.errorParser)
		});
	}

	// Get Processor Manifest
	// Endpoint: /get/manifest
	// Datalink string             `schema:"datalink,required"`
	// Manifest ProcessingManifest `schema:"-"`

	getManifest(datalink) {
		return new Promise((resolve, reject) => {
			request
				.post(`http://${this.serverAddress}/get/manifest`)
				.type('form')
				.send({ datalink })
				.then((res) => {
					resolve(JSON.parse(res.body.Data).Manifest);
				})
				.catch(this.errorParser)
		});
	}

	errorParser(err) {
		return new Promise((resolve, reject) => {
			const { Code, Data } = err.response.body;
			console.error(`${Code}: ${Data}`);
			reject(err);
		});
	}
};

export default WeatherClient;