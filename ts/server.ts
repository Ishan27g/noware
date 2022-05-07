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

app.post('/node/1', (req: Request, res: Response) => {
  
  /**
   * get context and actions for this request
   */

  let ctx = Context.get(req)
  let actions = Action.FromCtx(ctx)
  console.log("Request is noop ? ", Noop.ContainsNoop(ctx))
 
  console.log("Request action events -> ", actions?.Get())
  let event1 : Event = {name: 'server 1', meta: {'at':'server url 1 at NODE'} }
  actions?.Add(event1)

  if (actions === undefined || actions === null){
    res.send({'noop?':Noop.ContainsNoop(ctx)});
    return
  }
  res.send(JSON.stringify(actions?.Get()));
});

app.post('/node/2', (req: Request, res: Response) => {
  
  let ctx = Context.get(req)
  
  let actions = Action.FromCtx(ctx)
  let event1 : Event = {name: 'server 2', meta: {'at':'server url 2 at NODE'} }

  actions?.Add(event1)

  console.log("Request is noop ? ", Noop.ContainsNoop(ctx))
  console.log("Request action events -> ", actions?.Get())
 
  if (actions === undefined || actions === null){
    res.send({'noop?':Noop.ContainsNoop(ctx)});
    return
  }
  res.send(JSON.stringify(actions?.Get()));
});

app.listen(port, () => {
  console.log(`⚡️[server]: Server is running at https://localhost:${port}`);
});


