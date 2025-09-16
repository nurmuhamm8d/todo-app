export namespace main {
	
	export class CategoryDTO {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
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
	export class SubtaskDTO {
	    id: number;
	    taskId: number;
	    title: string;
	    completed: boolean;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new SubtaskDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.taskId = source["taskId"];
	        this.title = source["title"];
	        this.completed = source["completed"];
	        this.createdAt = source["createdAt"];
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
	    categoryId?: number;
	    tags?: string[];
	
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
	        this.categoryId = source["categoryId"];
	        this.tags = source["tags"];
	    }
	}

}

