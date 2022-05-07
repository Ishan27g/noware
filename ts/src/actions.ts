import e from "express"
import { type } from "os"

// export class Event {
//     name : string = ""
//     meta : object = {}
//     public metaString() : string{
//         return JSON.stringify(this.meta)
//     }
// }
export type Event = {
    name : string
    meta : object
}
export class Actions {
    private events : Event[] = []

    constructor(...events : Event[]){
        events.forEach(e => {
            this.Add(e)
        })
    }
    
    Add(e :Event){ 
        this.events.push(e)
    }
    Get(): Event[]{
        const events = this.events
        return events
    }
}   
