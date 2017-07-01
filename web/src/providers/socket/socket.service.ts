import {Injectable} from '@angular/core';

export enum OP { REG, SEND, ACK}
;

@Injectable()
export class SocketService {
  private ws: WebSocket = null;
  private url = '';
  private ackDict = new Dictionary<number, any>();
  private seqId = 0;

  constructor() {
  }

  start() {
    this.ws = new WebSocket(this.url);
    this.ws.binaryType = 'arraybuffer';
    // TODO send Function
    this.ws.onopen = () => {
      this.send();
    }
    this.ws.onclose = (error) => {
      console.error('ws error:', error);
    }
    this.ws.onclose = (res) => {
      console.log('-------ws close-------');
    }
  }

  private getSeq(buf: Uint8Array): number {
    return (buf[1] << 24) | (buf[2] << 16) | (buf[3] << 8) | (buf[4]);
  }

  private send(op: number, json?: string) {
    this.ackDict.setValue(++this.seqId, { op: op, json: json });
    this.sendWithSeq(op, this.seqId, json);
  }

  private sendWithSeq(op, seq: number, json?: string) {
    let buf: Uint8Array;
    let uintjson: Uint8Array;
    if (json) {
      console.debug(json);
      uintjson = this.stringToUint8(json);
      buf = new Uint8Array(uintjson.length + 5);
      for (var i = 5; i < buf.byteLength; i++) {
        buf[i] = uintjson[i - 5];
      }
    } else {
      buf = new Uint8Array(5);
    }

    //op
    buf[0] = 0xff & op;
    //seq
    buf[1] = 0xff & (seq >> 24);
    buf[2] = 0xff & (seq >> 16);
    buf[3] = 0xff & (seq >> 8);
    buf[4] = 0xff & seq;

    // this.waitForConnection(() => {
    this.ws.send(buf);
    // }, 1000);
  }

  ack(op: any, seq: number) {
    console.debug('op:%s seq:%d', op, seq);
    this.sendWithSeq(OP.ACK, seq);
  }

  private stringToUint8(str: string): Uint8Array {
    const bytes = new Array();
    let len, c;
    len = str.length;
    for (let i = 0; i < len; i++) {
      c = str.charCodeAt(i);
      if (c >= 0x010000 && c <= 0x10FFFF) {
        bytes.push(((c >> 18) & 0x07) | 0xF0);
        bytes.push(((c >> 12) & 0x3F) | 0x80);
        bytes.push(((c >> 6) & 0x3F) | 0x80);
        bytes.push((c & 0x3F) | 0x80);
      } else if (c >= 0x000800 && c <= 0x00FFFF) {
        bytes.push(((c >> 12) & 0x0F) | 0xE0);
        bytes.push(((c >> 6) & 0x3F) | 0x80);
        bytes.push((c & 0x3F) | 0x80);
      } else if (c >= 0x000080 && c <= 0x0007FF) {
        bytes.push(((c >> 6) & 0x1F) | 0xC0);
        bytes.push((c & 0x3F) | 0x80);
      } else {
        bytes.push(c & 0xFF);
      }
    }
    return new Uint8Array(bytes);
  }
}
