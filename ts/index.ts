import express, { Express, Request, Response } from 'express';
import dotenv from 'dotenv';
import Noop from './src/noop';
import Context from './src/context';
import { Action } from './src/action';
import { Event } from './src/actions'


dotenv.config();

const app: Express = express();
const port = process.env.PORT? process.env.PORT: 8082;

/**
 * noop/actions middleware 
 */
app.use(Noop.Middleware());

app.post('/1', (req: Request, res: Response) => {
  
  /**
   * get context and actions for this request
   */

  let ctx = Context.get(req)
  let actions = Action.FromCtx(ctx)
  console.log("Request is noop ? ", Noop.ContainsNoop(ctx))
  console.log("Request action events -> ", actions?.Get())
  let event = new Event()
  event.name = 'at server url 1'
  event.meta = {'at':'server url 1'}
  actions?.Add(event)
  if (actions === undefined || actions === null){
    res.send({'noop?':Noop.ContainsNoop(ctx)});
    return
  }
  res.send({'noop?':Noop.ContainsNoop(ctx), 'actions': JSON.stringify(actions?.Get())});
});

app.post('/2', (req: Request, res: Response) => {
  
  let ctx = Context.get(req)
  console.log(ctx)
  
  let actions = Action.FromCtx(ctx)
  console.log("actions 2", actions?.Get())
  let event = new Event()
  event.name = 'at server url 2'
  event.meta = {'at':'server url 2'}
  actions?.Add(event)

  console.log("Request is noop ? ", Noop.ContainsNoop(ctx))
  console.log("Request action events -> ", actions?.Get())
 
  if (actions === undefined || actions === null){
    res.send({'noop?':Noop.ContainsNoop(ctx)});
    return
  }
  res.send({'noop?':Noop.ContainsNoop(ctx), 'actions': JSON.stringify(actions?.Get())});
});

app.listen(port, () => {
  console.log(`⚡️[server]: Server is running at https://localhost:${port}`);
});


