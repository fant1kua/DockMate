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
	export class ImageInfo {
	    id: string;
	    name: string;
	    size: number;
	    tags: string[];
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
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

