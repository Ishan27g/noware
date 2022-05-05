import { IncomingHttpHeaders } from 'http'
import{ Headers } from 'node-fetch'
import { Actions, Event } from './actions'
import Context from './context'

export class  Action {
    private static actionKey = 'action-key'
    private static key() : string{
        return this.actionKey.toString()
    }

    static Header(from: Headers, action :Actions) : Headers{
        action.Get().forEach(e => {
            // todo : doesn't parse array as array
            from.append(Action.key(), "["+JSON.stringify(e)+"]")
        })
        return new Headers(from)!
    }
    static FromHeaders(from: IncomingHttpHeaders): Actions | null{
        const actionsString : string = from[Action.key()] as string // todo : doesn't parse array as array
        if (actionsString === undefined || actionsString === null){
            return null
        }
        return Action.from(actionsString)
    }
    private static from(actionsString: string) {
        let events: Event[] = []
        actionsString.split(", [").forEach(s => {
            s = s.replace("]", "")
            s = s.replace("[", "")
            events.push(JSON.parse(s.trim()))
        })
        return new Actions(...events)
    }

    static NewCtxWithActions(actions: Actions): Context|null{
        return newCtxWithActions(new Context(), Action.key(), actions)
    }
    static FromCtx(context :Context | null): Actions|null {
        if (context == null){
            return null
        }
        if (!containsActions(context, Action.key())){
            return null
        }
        let val = context.get(Action.key())
        if (val == null){
            return null
        }
        return val as Actions
    }
}
function newCtxWithActions(context: Context | null, key: string, action: Actions): Context {
    if (context == null) {
        context = new Context();
    }
    if (containsActions(context, key)) {
        return context;
    }
    context.set(key, action);
    return context;
}
function containsActions(context: Context | null, key: string): boolean {
    if (context == null) {
        return false;
    }
    let val = context.get(key);
    if (val == null) {
        return false;
    }
    return true;
}
