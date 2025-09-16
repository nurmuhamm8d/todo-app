export namespace main {
	
	export class StatsDTO {
	    total: number;
	    active: number;
	    completed: number;
	    overdue: number;
	
	    static createFrom(source: any = {}) {
	        return new StatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.active = source["active"];
	        this.completed = source["completed"];
	        this.overdue = source["overdue"];
	    }
	}
	export class TaskDTO {
	    id: number;
	    title: string;
	    priority: string;
	    completed: boolean;
	    createdAt: string;
	    completedAt?: string;
	    dueDate?: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.priority = source["priority"];
	        this.completed = source["completed"];
	        this.createdAt = source["createdAt"];
	        this.completedAt = source["completedAt"];
	        this.dueDate = source["dueDate"];
	    }
	}

}

