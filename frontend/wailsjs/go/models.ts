export namespace main {
	
	export class DfaLoadResult {
	    filename: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new DfaLoadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.content = source["content"];
	    }
	}
	export class InLoadResult {
	    filename: string;
	    inputLines: string[];
	
	    static createFrom(source: any = {}) {
	        return new InLoadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.inputLines = source["inputLines"];
	    }
	}

}

