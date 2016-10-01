extern crate ws;
extern crate env_logger;

use ws::{listen, Sender, Factory, Handler};
use std::collections::HashMap;
use std::collections::hash_map::Entry;
use std::rc::Rc;
use std::cell::{RefCell, Ref};

fn main() {

    let mut senders_by_thread = HashMap::<usize, Rc<RefCell<Vec<ws::Sender>>>>::new();
    // Setup logging
    env_logger::init().unwrap();

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

struct MyHandler {
    ws: Sender,
    is_server: bool,
    senders: Rc<RefCell<Vec<ws::Sender>>>,
}

impl Handler for MyHandler {
    fn on_message(&mut self, msg: ws::Message) -> ws::Result<()> {
        println!("Server got message '{}'. ", msg);
        for sender in self.senders.borrow().iter() {
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
    senders_by_thread: HashMap<usize, Rc<RefCell<Vec<ws::Sender>>>>,
}

impl Factory for MyFactory {
    type Handler = MyHandler;

    fn connection_made(&mut self, ws: Sender) -> MyHandler {
        MyHandler {
            ws: ws,
            // default to client
            is_server: false,
            senders: Rc::new(RefCell::new(vec![])),
        }
    }

    fn server_connected(&mut self, ws: Sender) -> MyHandler {
        match self.senders_by_thread.entry(0) {
            Entry::Occupied(o) => {
                o.into_mut().borrow_mut().push(ws.clone());
            }
            Entry::Vacant(v) => {
                v.insert(Rc::new(RefCell::new(vec![ws.clone()])));
            }
        }

        MyHandler {
            ws: ws,
            is_server: true,
            senders: self.senders_by_thread.get(&0).unwrap().clone(), // should always have senders since if it doesn't exists it will be added
        }
    }
}
