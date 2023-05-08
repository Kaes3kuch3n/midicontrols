export namespace gui {
	
	export class Settings {
	    actions: {[key: number]: };
	    selectedDevice: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.actions = source["actions"];
	        this.selectedDevice = source["selectedDevice"];
	    }
	}

}

export namespace midihandler {
	
	export class Event {
	    type: number;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.value = source["value"];
	    }
	}

}

