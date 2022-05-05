import { Request } from 'express';

export default class Context {
    // context binding for a request
    private static _bindings = new WeakMap<Request, Context>();

    // context values
    private _values = new Map<String, Object>()
    
    constructor () {}
    
    // binds the context to this request
    static bind (req: Request, ctx : Context | null) : void {
        ctx == null ? Context._bindings.set(req, new Context()): Context._bindings.set(req, ctx);
    }
    // returns the context or null for this request
    static get (req: Request) : Context | null {
        return Context._bindings.get(req) || null;
    }
    // sets a value for this context
    public set (key: String, val: Object) : void {
        this._values.set(key, val);
    }
    // sets a value for this context
    public get (key: String) : Object | null {
        return this._values.get(key) || null;
    }
    
}