export namespace main {
	
	export class Log {
	    time: number;
	    type: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Log(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.type = source["type"];
	        this.message = source["message"];
	    }
	}
	export class Proxy {
	    id: number;
	    name: string;
	    create_time: number;
	    s: string;
	    p: number;
	    b: string;
	    l: number;
	    k: string;
	    m: string;
	    o: string;
	    op: string;
	    oo: string;
	    oop: string;
	    t: number;
	    f: string;
	    status: number;
	    run_status: string;
	
	    static createFrom(source: any = {}) {
	        return new Proxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.create_time = source["create_time"];
	        this.s = source["s"];
	        this.p = source["p"];
	        this.b = source["b"];
	        this.l = source["l"];
	        this.k = source["k"];
	        this.m = source["m"];
	        this.o = source["o"];
	        this.op = source["op"];
	        this.oo = source["oo"];
	        this.oop = source["oop"];
	        this.t = source["t"];
	        this.f = source["f"];
	        this.status = source["status"];
	        this.run_status = source["run_status"];
	    }
	}

}

