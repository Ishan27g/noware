
export class Event {
    name : string = ""
    meta : object = {}
    public metaString() : string{
        return JSON.stringify(this.meta)
    }
}
export class Actions {
    private events : Event[] = []

    constructor(...events : Event[]){
        this.events = events
    }
    
    Add(e :Event){ this.events.push(e) }
    Get(): Event[]{
        const events = this.events
        return events
    }
}   
