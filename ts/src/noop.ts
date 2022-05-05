import  { Request } from 'express';

import Context from "./context";
import{ Headers } from 'node-fetch'
import { IncomingHttpHeaders } from "http";
import { Action } from './action';
import { Actions, Event } from './actions'

export default class Noop {
    private static noopKey = "noop-key"

    static Middleware(): any {
        return (req: Request, res: any, next: any) => {
            let ctx = Context.get(req);
            // ctx is always null since it uses a weakmap[key = request object] => new key per request
            // irrespective, weakmap[key] would anyways be garbage collected at the end of scope for which the key was added
            if (ctx == null) {  // todo : remove null check
                if (Noop.checkHeader(req.headers)) {
                    console.log("has noop header")
                    let actions = Action.FromHeaders(req.headers)
                    if (actions != null){
                        ctx = Action.NewCtxWithActions(actions)
                    }
                    ctx = Noop.NewCtxWithNoop(ctx);
                }
                Context.bind(req, ctx);
            }
            next();
        };
    }
    private static key() : string{
        return this.noopKey.toString()
    }
    static MakeHeader(from: any, actions?: Actions | null) : Headers{
        if (actions === undefined || actions === null){
            return Noop.header(from)
        }
        return Action.Header(Noop.header(from), actions)
    }
    private static header(from: any) : Headers{
        from[this.key()] = 'true'
        return new Headers(from)!
    }
    private static checkHeader(h : IncomingHttpHeaders) : boolean{
        let hasNoop = h[this.key()]
        if (hasNoop === undefined || hasNoop === null){
            return false
        }
        return hasNoop === 'true' ? true : false
    }
    static ContainsNoop(context:Context | null) : boolean{
        return containsNoop(context, this.key())
    }

    static NewCtxWithNoop(context: Context| null) : Context{
       return newCtxWithNoop(context, true, this.key())
    }
}

function containsNoop(context:Context | null, key : string) : boolean{
    if (context == null){
        return false
    }
    let val = context.get(key)
    if (val == null){
        return false
    }
    return val.toString() == 'true'
}

function newCtxWithNoop(context: Context| null, isNoop:boolean, key : string) : Context{
    if (context == null){
        context = new Context()
    }
    if (containsNoop(context,key)){
        return context
    }
    context.set(key,isNoop)
    return context
}