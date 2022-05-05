import fetch from 'node-fetch';
import { Actions,Event } from './src/actions';
import Noop from './src/noop';

const existingHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'Access-Control-Allow-Headers': 'Content-Type',
  }

let someData = {name: "anyone"};

/**
 * Add action events
 */
let event1 = new Event()
event1.name = 'event 1'
event1.meta = {'from':'client 1'}
let actions = new Actions(event1)

fetch('http://localhost:8082/1', {
    method: 'POST',
    body: JSON.stringify(someData),
    /**
     * update existing headers with noop/action headers
     */
    headers: Noop.MakeHeader(existingHeaders, actions),
})
.then((result: { json: () => any; }) => result.json())
.then((jsonformat: any)=>{
  console.log(jsonformat)
  let e: Event[] = JSON.parse(jsonformat.actions)
  let newActions : Actions = new Actions(...e)
  let newEvent = new Event()
  newEvent.name = 'event 2'
  newEvent.meta = {'from':'client 2'}
  newActions.Add(newEvent)

  fetch('http://localhost:8082/2', {
    method: 'POST',
    headers: Noop.MakeHeader(existingHeaders, newActions),
  })
  .then((result: { json: () => any; }) => result.json())
  .then((jsonformat: any)=>{
    console.log(jsonformat)
    if (jsonformat.actions !== undefined){
      let e: Event[] = JSON.parse(jsonformat.actions)
      let newActions : Actions = new Actions(...e)
      console.log("actions events get ",newActions.Get())
    }
  });

});

