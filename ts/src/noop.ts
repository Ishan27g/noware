import  { Request } from 'express';

import Context from "./context";
import{ Headers } from 'node-fetch'
import { IncomingHttpHeaders } from "http";
import { Action, Actions } from './actions';

export default class Noop {
    private static noopKey = "noop-key"

    static Middleware(): any {
        return (req: Request, res: any, next: any) => {
            let ctx = Context.get(req);
            if (ctx == null) { // todo: is the ctx is always null ?
                if (Noop.CheckHeader(req.headers)) {
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
    static HeaderKey() : string{
        return this.noopKey.toString()
    }
    static Header(from: any) : Headers{
        from[this.HeaderKey()] = 'true'
        return new Headers(from)!
    }
    static CheckHeader(h : IncomingHttpHeaders) : boolean{
        let hasNoop = h[this.HeaderKey()]
        if (hasNoop === undefined || hasNoop === null){
            return false
        }
        return hasNoop === 'true' ? true : false
    }
    static ContainsNoop(context:Context | null) : boolean{
        return containsNoop(context, this.HeaderKey())
    }

    static NewCtxWithNoop(context: Context| null) : Context{
       return newCtxWithNoop(context, true, this.HeaderKey())
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