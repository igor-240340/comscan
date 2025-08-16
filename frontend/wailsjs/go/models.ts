export namespace main {
	
	export class ComPortInfo {
	    Name: string;
	    Usb: string;
	    Vid: string;
	    Pid: string;
	    SentData: string;
	    ReceivedData: string;
	
	    static createFrom(source: any = {}) {
	        return new ComPortInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Usb = source["Usb"];
	        this.Vid = source["Vid"];
	        this.Pid = source["Pid"];
	        this.SentData = source["SentData"];
	        this.ReceivedData = source["ReceivedData"];
	    }
	}

}

