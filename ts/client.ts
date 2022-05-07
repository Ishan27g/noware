import fetch from 'node-fetch';
import { Actions,Event } from './src/actions';
import Noop from './src/noop';

const existingHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'Access-Control-Allow-Headers': 'Content-Type',
  }

const someData = {name: "anyone"};

/**
 * Add action events
 */
let event1 : Event = {name: 'event 1', meta: {'at':'1.TS@client'} }
let event2 : Event = {name: 'event 2', meta: {'at':'2.TS@client'} }

const actions = new Actions(event1, event2)

const headersWithNoopAndActions = Noop.MakeHeader(existingHeaders, actions)

fetch('http://localhost:8082/node/1', {
    method: 'POST',
    body: JSON.stringify(someData),
    /**
     * update existing headers with noop/action headers
     */
    headers: headersWithNoopAndActions,
})
.then((result: { json: () => any; }) => result.json())
.then((jsonformat: any)=>{
  console.log(jsonformat)
  let e: Event[] = JSON.parse(JSON.stringify(jsonformat))
  console.log("actions events response 1 ",e)
  let newActions : Actions = new Actions(...e)
  let event2 : Event = {name: 'event 3', meta: {'at':'TS@client to NODE'} }
  let event22 : Event = {name: 'event 4', meta: {'at':'TS@client to NODE'} }

  newActions.Add(event2)
  newActions.Add(event22)

  fetch('http://localhost:8082/node/2', {
    body: JSON.stringify(someData),
    method: 'POST',
    headers: Noop.MakeHeader(existingHeaders, newActions),
  })
  .then((result: { json: () => any; }) => result.json())
  .then((jsonformat: any)=>{
      console.log(jsonformat)
      let e: Event[] = JSON.parse(JSON.stringify(jsonformat))
      let newActions : Actions = new Actions(...e)
      console.log("actions events response 2 ",newActions.Get())
  });

});
// same as above with go backend
setTimeout(()=>{
fetch('http://localhost:8081/go/1', {
    body: JSON.stringify(someData),
    method: 'POST',
    headers: headersWithNoopAndActions,
  })
  .then((result: { json: () => any; }) => result.json())
  .then((jsonformat: any)=>{
      let e: Event[] = JSON.parse(JSON.stringify(jsonformat))
      let newActions : Actions = new Actions(...e)
      console.log("actions events response 1 ",newActions.Get())

      let event2 : Event = {name: 'event 3', meta: {'at':'TS@client to GO'} }
      let event22 : Event = {name: 'event 4', meta: {'at':'TS@client to GO'} }
    
      newActions.Add(event2)
      newActions.Add(event22)

      fetch('http://localhost:8081/go/2', {
        body: JSON.stringify(someData),
        method: 'POST',
        headers: Noop.MakeHeader(existingHeaders, newActions),
      })
      .then((result: { json: () => any; }) => result.json())
      .then((jsonformat: any)=>{
          let e: Event[] = JSON.parse(JSON.stringify(jsonformat))
          let newActions : Actions = new Actions(...e)
          console.log("actions events response 2 ",newActions.Get())
      });
  });
}, 1000)
