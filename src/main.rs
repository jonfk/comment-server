extern crate ws;
extern crate env_logger;

use ws::{listen, Sender, Factory, Handler};
use std::collections::HashMap;
use std::collections::hash_map::Entry;
use std::rc::Rc;
use std::cell::RefCell;

fn main() {

    let mut senders_by_thread = HashMap::<usize, RefCell<Vec<ws::Sender>>>::new();
    // Setup logging
    env_logger::init().unwrap();

    // // Listen on an address and call the closure for each connection
    // if let Err(error) = listen("127.0.0.1:3012", |out| {
    //     // subscribe
    //     // match senders_by_thread.entry(0) {
    //     //     Entry::Occupied(o) => {
    //     //         o.into_mut().push(out.clone());
    //     //     }
    //     //     Entry::Vacant(v) => {
    //     //         v.insert(vec![out.clone()]);
    //     //     }
    //     // };

    //     // The handler needs to take ownership of out, so we use move
    //     move |msg: ws::Message| {
    //         // publish
    //         // if let Some(senders) = senders_by_thread.get(&0) {
    //         //     for sender in senders {
    //         //         sender.send(msg.clone());
    //         //     }
    //         // }

    //         // Handle messages received on this connection
    //         println!("Server got message '{}'. ", msg);

    //         // Use the out channel to send messages back
    //         out.send(msg)
    //     }
    // }) {
    //     // Inform the user of failure
    //     println!("Failed to create WebSocket due to {:?}", error);
    // }
    let factory = MyFactory { senders_by_thread: senders_by_thread };
    match ws::WebSocket::new(factory) {
        Ok(ws) => {
            if let Err(error) = ws.listen("127.0.0.1:3012") {
                // Inform the user of failure
                println!("Failed to create WebSocket due to {:?}", error);
            }
        }
        Err(error) => {
            println!("Failed to create WebSocket due to {:?}", error);
        }
    }
}

struct MyHandler<'a> {
    ws: Sender,
    is_server: bool,
    senders: &'a RefCell<Vec<ws::Sender>>,
}

impl<'a> Handler for MyHandler<'a> {
    fn on_message(&mut self, msg: ws::Message) -> ws::Result<()> {
        println!("Server got message '{}'. ", msg);
        for sender in self.senders.borrow() {
            sender.send(msg.clone());
        }
        self.ws.send(msg)
    }

    fn on_close(&mut self, code: ws::CloseCode, reason: &str) {
        println!("WebSocket closing for ({:?}) {}", code, reason);
        // println!("Shutting down server after first connection closes.");
        // self.ws.shutdown().unwrap();
    }
}

struct MyFactory {
    senders_by_thread: HashMap<usize, RefCell<Vec<ws::Sender>>>,
}

impl Factory for MyFactory<'a> {
    type Handler = MyHandler<'a>;

    fn connection_made(&mut self, ws: Sender) -> MyHandler {
        MyHandler {
            ws: ws,
            // default to client
            is_server: false,
            senders:RefCell::new(vec![]),
        }
    }

    fn server_connected(&mut self, ws: Sender) -> MyHandler {
        match self.senders_by_thread.entry(0) {
            Entry::Occupied(o) => {
                o.into_mut().borrow_mut().push(ws.clone());
            }
            Entry::Vacant(v) => {
                v.insert(RefCell::new(vec![ws.clone()]));
            }
        }

        MyHandler {
            ws: ws,
            is_server: true,
            senders: self.senders_by_thread.get(&0).unwrap(), // should always have senders since if it doesn't exists it will be added
        }
    }
}
