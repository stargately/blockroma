import EventEmitter from "events";

const fetcherEmitter = new EventEmitter();

export const publisher = {
  async publish(channel: string, message: string): Promise<void> {
    fetcherEmitter.emit(channel, message);
  },
};

export const subscriber = {
  async subscribe(
    channel: string,
    cb: (message: string) => Promise<void>
  ): Promise<void> {
    fetcherEmitter.on(channel, cb);
  },
};
