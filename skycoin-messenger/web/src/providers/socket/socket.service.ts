import { Injectable } from '@angular/core';
import { ImHistoryMessage, RecentItem, UserInfo } from './msg';
import { Subject } from 'rxjs/Subject';
import * as Collections from 'typescript-collections';
import { UserService } from '../user/user.service';
import { environment } from '../../environments/environment';

export enum OP { REG, SEND, ACK };
export enum PUSH { ACK, REG, MSG };

@Injectable()
export class SocketService {
  private ws: WebSocket = null;
  private url = 'ws://localhost:8082/ws';
  // private url = 'ws://192.168.33.104:8082/ws';
  // private ackDict = new Dictionary<number, any>();
  key = ''
  chattingUser = '';
  private seqId = 0;
  recent_list: Array<RecentItem> = [];
  // private history = new Collections.LinkedDictionary<string, Array<ImHistoryMessage>>()
  private historySubject = new Subject<Map<string, Collections.LinkedList<ImHistoryMessage>>>();
  histories = new Map<string, Collections.LinkedList<ImHistoryMessage>>();
  userInfo = new Map<string, UserInfo>();
  chatHistorys = this.historySubject.asObservable();
  constructor(private user: UserService) {
    if (environment.production) {
      this.url = 'ws://yiqishare.com:8082/ws';
    }
    this.historySubject.subscribe((data: Map<string, Collections.LinkedList<ImHistoryMessage>>) => {
      this.histories = data;
    })
  }
  getRencentListIndex(key: string) {
    return this.recent_list.findIndex(v => v.name === key);
  }
  addHint(key, msg: string) {
    const index = this.getRencentListIndex(key);
    if (index <= -1) {
      console.log('add new item', key);
      const icon = this.user.getRandomMatch();
      this.recent_list.push({ name: key, last: msg, unRead: 1, icon: icon });
      this.userInfo.set(key, { Icon: icon });
    } else {
      if (key !== this.chattingUser) {
        console.log('add un read:', this.recent_list[index]);
        this.recent_list[index].unRead += 1;
      }
      // tslint:disable-next-line:no-unused-expression
      this.recent_list[index].last = msg;
    }
  }
  start() {
    this.ws = new WebSocket(this.url);
    this.ws.binaryType = 'arraybuffer';
    this.ws.onopen = () => {
      this.send(OP.REG, JSON.stringify({ Address: 'localhost:8080' }));
      // this.send(OP.REG);
    }
    this.ws.onmessage = (event) => {
      this.handle(event.data);
    }
    this.ws.onerror = (error) => {
      console.error('-----ws error-------', error);
    }
    this.ws.onclose = (res) => {
      console.log('-------ws close-------', res);
    }
  }
  private getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min);
  }
  msg(chattingKey, message: string) {
    this.send(OP.SEND, JSON.stringify({ PublicKey: chattingKey, Msg: message }));
  }

  private handle(data: ArrayBuffer) {
    const buf = new Uint8Array(data);
    // console.log('orgin data:', this.toHexString(buf));
    const op = buf[0]
    const metaData = this.utf8ArrayToStr(buf.slice(5));
    let json = null;
    if (metaData) {
      json = JSON.parse(metaData);
    }
    switch (op) {
      case PUSH.ACK:
        this.ack(op, this.getSeq(buf));
        break;
      case PUSH.REG:
        this.key = json.PublicKey;
        this.userInfo.set(this.key, { Icon: this.user.getRandomMatch() });
        break;
      case PUSH.MSG:
        const now = new Date().getTime();
        let list = this.histories.get(json.From)
        if (list === undefined) {
          list = new Collections.LinkedList<ImHistoryMessage>();
        }
        list.add({ From: json.From, Msg: json.Msg, Timestamp: now }, 0);
        this.addHint(json.From, json.Msg);
        this.saveHistorys(json.From, list);
        this.ack(op, this.getSeq(buf));
        break;
    }
  }
  saveHistorys(key: string, msgList: Collections.LinkedList<ImHistoryMessage>) {
    this.histories.set(key, msgList);
    this.historySubject.next(this.histories);
  }

  private toHexString(byteArray) {
    return Array.from(byteArray, (byte: number) => {
      return ('0' + (byte & 0xFF).toString(16)).slice(-2);
    }).join('')
  }
  private getSeq(buf: Uint8Array): number {
    return (buf[1] << 24) | (buf[2] << 16) | (buf[3] << 8) | (buf[4]);
  }

  private send(op: number, json?: string) {
    // this.ackDict.setValue(++this.seqId, { op: op, json: json });
    this.sendWithSeq(op, ++this.seqId, json);
  }

  private sendWithSeq(op, seq: number, json?: string) {
    let buf: Uint8Array;
    let uintjson: Uint8Array;
    if (json) {
      // console.log('send json:', json);
      // console.log('send seq:', seq);
      uintjson = this.stringToUint8(json);
      buf = new Uint8Array(uintjson.length + 5);
      for (let i = 5; i < buf.byteLength; i++) {
        buf[i] = uintjson[i - 5];
      }
    } else {
      buf = new Uint8Array(5);
    }

    // op
    buf[0] = 0xff & op;
    // seq
    buf[1] = 0xff & (seq >> 24);
    buf[2] = 0xff & (seq >> 16);
    buf[3] = 0xff & (seq >> 8);
    buf[4] = 0xff & seq;

    // this.waitForConnection(() => {
    this.ws.send(buf);
    // }, 1000);
  }

  ack(op: any, seq: number) {
    // console.log('op:%s seq:%d', op, seq);
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
  private utf8ArrayToStr(array) {
    let out, i, len, c;
    let char2, char3;

    out = '';
    len = array.length;
    i = 0;
    while (i < len) {
      c = array[i++];
      switch (c >> 4) {
        case 0:
        case 1:
        case 2:
        case 3:
        case 4:
        case 5:
        case 6:
        case 7:
          // 0xxxxxxx
          out += String.fromCharCode(c);
          break;
        case 12:
        case 13:
          // 110x xxxx   10xx xxxx
          char2 = array[i++];
          out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));
          break;
        case 14:
          // 1110 xxxx  10xx xxxx  10xx xxxx
          char2 = array[i++];
          char3 = array[i++];
          out += String.fromCharCode(((c & 0x0F) << 12) |
            ((char2 & 0x3F) << 6) |
            ((char3 & 0x3F) << 0));
          break;
      }
    }

    return out;
  }
}
