import fetch from 'node-fetch';
import { json } from 'stream/consumers';
import { Action, Actions, Event } from './src/actions';
import Noop from './src/noop';

const someHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'Access-Control-Allow-Headers': 'Content-Type',
  }

let someData = {userId: "01234"};

let headerWithNoop = Noop.Header(someHeaders)
let event1 = new Event()
event1.name = 'event 1'
event1.meta = {'from':'client 1'}
let actions = new Actions(event1)

fetch('http://localhost:8000/1', {
    method: 'POST',
    body: JSON.stringify(someData),
    headers: Action.Header(headerWithNoop, actions),
})
.then((result: { json: () => any; }) => result.json())
.then((jsonformat: any)=>{

  let e: Event[] = JSON.parse(jsonformat.actions)
  let newActions : Actions = new Actions(...e)
  let newEvent = new Event()
  newEvent.name = 'event 2'
  newEvent.meta = {'from':'client 2'}
  newActions.Add(newEvent)

  fetch('http://localhost:8000/2', {
    method: 'POST',
    headers: Action.Header(Noop.Header(someHeaders), newActions),
  })
  .then((result: { json: () => any; }) => result.json())
  .then((jsonformat: any)=>{

    let e: Event[] = JSON.parse(jsonformat.actions)
    let newActions : Actions = new Actions(...e)
    console.log("actions events get ",newActions.Get())
  });

});

