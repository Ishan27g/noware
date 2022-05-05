import express, { Express, Request, Response } from 'express';
import dotenv from 'dotenv';
import Noop from './src/noop';
import Context from './src/context';
import { Action, Actions, Event } from './src/actions';
import { json } from 'stream/consumers';

dotenv.config();

const app: Express = express();
const port = process.env.PORT? process.env.PORT: 8000;

// noop and actions middleware 
app.use(Noop.Middleware());

app.post('/1', (req: Request, res: Response) => {
  
  let ctx = Context.get(req)
  let actions = Action.FromHeaders(req.headers)

  console.log("Request has noop ? ", Noop.ContainsNoop(ctx))
  console.log("Request action events ? ", actions?.Get())
 
  let event = new Event()
  event.name = 'at server url 1'
  event.meta = {'at':'server url 1'}
  actions?.Add(event)
  
  res.send({'noop?':Noop.ContainsNoop(ctx), 'actions': JSON.stringify(actions?.Get())});
  // res.send({'noop?':Noop.ContainsNoop(ctx), 'actions': JSON.parse(JSON.stringify(actions?.Get()))});
});

app.post('/2', (req: Request, res: Response) => {
  
  let ctx = Context.get(req)
  let actions = Action.FromHeaders(req.headers)

  let event = new Event()
  event.name = 'at server url 2'
  event.meta = {'at':'server url 2'}
  actions?.Add(event)

  console.log("Request has noop ? ", Noop.ContainsNoop(ctx))
  console.log("Request action events ? ", actions?.Get())
 
  res.send({'noop?':Noop.ContainsNoop(ctx), 'actions': JSON.stringify(actions?.Get())});
  // res.send({'noop?':Noop.ContainsNoop(ctx), 'actions': JSON.parse(JSON.stringify(actions?.Get()))});
});

app.listen(port, () => {
  console.log(`⚡️[server]: Server is running at https://localhost:${port}`);
});


