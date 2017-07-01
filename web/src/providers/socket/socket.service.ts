import { Injectable } from '@angular/core';

@Injectable()
export class SocketService {
  private ws: WebSocket = null;
  private url = '';
  constructor() {
  }

  /**
   * Register
   */
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

  send(op?: number, json?: string) {
    const buf: Uint8Array = new Uint8Array(13);
    let uintjson: Uint8Array;
    if (json) {
      uintjson = this.stringToUint8(json);
    }
    buf[0] = 0xff & (op >> 8);
    buf[1] = 0xff & op;
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
