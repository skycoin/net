import { Component, OnInit, ViewChild, ElementRef, OnDestroy } from '@angular/core';
import { ApiService } from '../../service/api/api.service';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { MatDialogRef } from '@angular/material';
import * as Terminal from 'xterm';
import 'rxjs/add/observable/interval';

@Component({
  selector: 'app-terminal',
  templateUrl: 'terminal.component.html',
  styleUrls: ['./terminal.component.scss']
})

export class TerminalComponent implements OnInit, OnDestroy {
  @ViewChild('edit') edit: ElementRef;
  @ViewChild('container') conainer: ElementRef;
  index = -1;
  addr = '';
  task: Subscription = null;
  xterm: Terminal = null;
  ws: WebSocket = null;
  url = 'ws://127.0.0.1:8000/term';
  isWrite = true;
  constructor(
    private api: ApiService,
    private dialogRef: MatDialogRef<TerminalComponent>
  ) { }

  ngOnInit() {
    this.ws = new WebSocket(`${this.url}?url=ws://${this.addr}/node/run/term`);
    this.ws.binaryType = 'arraybuffer';
    this.ws.onopen = (ev) => {
      this.start();
    };
    this.ws.onclose = (ev: CloseEvent) => {
      this.isWrite = false;
      this.xterm.writeln('Connection interrupted...');
    };
    this.ws.onerror = (ev: Event) => {
      this.isWrite = false;
      this.xterm.writeln('Connection interrupted...');
    };
  }
  ngOnDestroy() {
    this.ws.close();
  }
  send(data) {
    if (this.isWrite) {
      this.ws.send(this.stringToUint8(data));
    }
  }
  start() {
    this.xterm = new Terminal({
      cursorBlink: true,
    });
    this.xterm.open(this.conainer.nativeElement);
    this.ws.onmessage = (evt) => {
      if (evt.data instanceof ArrayBuffer) {
        this.xterm.write(this.utf8ArrayToStr(new Uint8Array(evt.data)));
      } else {
        console.log('ws data:', evt.data);
      }
    };
    this.xterm.on('data', (data) => {
      this.send(data);
    });
  }
  _close(ev?: Event) {
    if (ev) {
      ev.stopImmediatePropagation();
      ev.stopPropagation();
      ev.preventDefault();
    }
    this.dialogRef.close();
  }
  utf8ArrayToStr(array) {
    let out, i, len, c;
    let char2, char3;

    out = '';
    len = array.length;
    i = 0;
    while (i < len) {
      c = array[i++];
      // tslint:disable-next-line:no-bitwise
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
          // tslint:disable-next-line:no-bitwise
          out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));
          break;
        case 14:
          // 1110 xxxx  10xx xxxx  10xx xxxx
          char2 = array[i++];
          char3 = array[i++];
          // tslint:disable-next-line:no-bitwise
          out += String.fromCharCode(((c & 0x0F) << 12) |
            // tslint:disable-next-line:no-bitwise
            ((char2 & 0x3F) << 6) |
            // tslint:disable-next-line:no-bitwise
            ((char3 & 0x3F) << 0));
          break;
      }
    }

    return out;
  }
  stringToUint8(str: string): Uint8Array {
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
