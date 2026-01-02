export namespace yaml {
	
	export class Vector2 {
	    X: number;
	    Y: number;
	
	    static createFrom(source: any = {}) {
	        return new Vector2(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	    }
	}
	export class UIControl {
	    Class: string;
	    CustomClass: string;
	    Name: string;
	    Position?: Vector2;
	    Size?: Vector2;
	    Pivot?: Vector2;
	    Visible?: boolean;
	    Input?: boolean;
	    Classes: string;
	    Prototype: string;
	    Components: Record<string, any>;
	    Children: UIControl[];
	    Properties: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new UIControl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Class = source["Class"];
	        this.CustomClass = source["CustomClass"];
	        this.Name = source["Name"];
	        this.Position = this.convertValues(source["Position"], Vector2);
	        this.Size = this.convertValues(source["Size"], Vector2);
	        this.Pivot = this.convertValues(source["Pivot"], Vector2);
	        this.Visible = source["Visible"];
	        this.Input = source["Input"];
	        this.Classes = source["Classes"];
	        this.Prototype = source["Prototype"];
	        this.Components = source["Components"];
	        this.Children = this.convertValues(source["Children"], UIControl);
	        this.Properties = source["Properties"];
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
	export class Header {
	    Version: number;
	
	    static createFrom(source: any = {}) {
	        return new Header(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Version = source["Version"];
	    }
	}
	export class UIPackage {
	    Header: Header;
	    ImportedPackages: string[];
	    Prototypes: UIControl[];
	    ExternalPackages: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new UIPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Header = this.convertValues(source["Header"], Header);
	        this.ImportedPackages = source["ImportedPackages"];
	        this.Prototypes = this.convertValues(source["Prototypes"], UIControl);
	        this.ExternalPackages = source["ExternalPackages"];
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
	export class FileData {
	    path: string;
	    content: string;
	    package?: UIPackage;
	    assets: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.content = source["content"];
	        this.package = this.convertValues(source["package"], UIPackage);
	        this.assets = source["assets"];
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
	
	
	

}

