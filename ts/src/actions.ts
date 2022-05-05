import e from 'express'
import { IncomingHttpHeaders } from 'http'
import{ Headers } from 'node-fetch'
import Context from './context'

export class Event {
    name : string = ""
    meta : object = {}
    public metaString() : string{
        return JSON.stringify(this.meta)
    }
}
export class  Action {
    private static actionKey = 'action-key'
    static HeaderKey() : string{
        return this.actionKey.toString()
    }

    static Header(from: Headers, action :Actions) : Headers{
        let val = ""
        let copy = from
        action.Get().forEach(e => {
            from.append(Action.HeaderKey(), "["+JSON.stringify(e)+"]")            
        })
        val = copy.get(Action.HeaderKey())!
        from.set(Action.HeaderKey(), val)
        return new Headers(from)!
    }
    static FromHeaders(from: IncomingHttpHeaders): Actions | null{
        const hasActions : string = from[Action.HeaderKey()] as string
        if (hasActions === undefined || hasActions === null){
            return null
        }
        return Action.From(hasActions)
    }
    static From(hasActions: string) {
        let events: Event[] = []
        hasActions.split(", [").forEach(s => {
            s = s.trim()
            s = s.replace("]", "")
            s = s.replace("[", "")
            let e = JSON.parse(s)
            events.push(e)
        })
        return new Actions(...events)
    }

    static NewCtxWithActions(actions: Actions): Context|null{
        return newCtxWithActions(new Context(), Action.HeaderKey(), actions)
    }
    static FromCtx(ctx :Context): Actions|null {
        return fromCtx(ctx, Action.HeaderKey())
    }
}
export class Actions {
    private events : Event[] = []

    constructor(...events : Event[]){
        this.events = events
    }
    
    Add(e :Event){
        this.events.push(e)
    }
    Get(): Event[]{
        const events = this.events
        return events
    }
    

}   
function fromCtx(context: Context, key : string) :Actions|null{
    if (context == null){
        return null
    }
    if (!containsActions(context, key)){
        return null
    }
    let val = context.get(key)
    if (val == null){
        return null
    }
    let actions = new Actions(...JSON.parse(JSON.stringify(val)))
    return actions

}
function newCtxWithActions(context: Context| null, key : string, action : Actions) : Context{
    if (context == null){
        context = new Context()
    }
    if (containsActions(context,key)){
        return context
    }
    context.set(key,action)
    return context
}
function containsActions(context:Context | null, key : string) : boolean{
    if (context == null){
        return false
    }
    let val = context.get(key)
    if (val == null){
        return false
    }
    return true
}