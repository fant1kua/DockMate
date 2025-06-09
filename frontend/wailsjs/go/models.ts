export namespace app {
	
	export class ContainerInfo {
	    id: string;
	    names: string[];
	    image: string;
	    status: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new ContainerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.names = source["names"];
	        this.image = source["image"];
	        this.status = source["status"];
	        this.state = source["state"];
	    }
	}
	export class ContainersGroup {
	    name: string;
	    containers: ContainerInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ContainersGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.containers = this.convertValues(source["containers"], ContainerInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImageInfo {
	    id: string;
	    size: number;
	    tags: string[];
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.size = source["size"];
	        this.tags = source["tags"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class NetworkInfo {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class VolumeInfo {
	    id: string;
	    name: string;
	    size: number;
	    tags: string[];
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new VolumeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.tags = source["tags"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

