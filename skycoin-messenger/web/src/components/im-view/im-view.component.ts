import {
  Component,
  OnInit,
  ViewEncapsulation,
  Input,
  SimpleChanges,
  ViewChild,
  AfterViewChecked,
  OnChanges
} from '@angular/core';
import { SocketService } from '../../providers';
import { ImHistoryMessage } from '../../providers';
import * as Collections from 'typescript-collections';
import { ImHistoryViewComponent } from '../im-history-view/im-history-view.component';

@Component({
  selector: 'app-im-view',
  templateUrl: './im-view.component.html',
  styleUrls: ['./im-view.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImViewComponent implements OnInit, OnChanges {
  chatList: Collections.LinkedList<ImHistoryMessage>;
  msg = '';
  @ViewChild(ImHistoryViewComponent) historyView: ImHistoryViewComponent;
  @Input() chatting = '';
  constructor(public socket: SocketService) { }

  ngOnInit() {
  }
  ngOnChanges(changes: SimpleChanges) {
    const previousValue = changes.chatting.previousValue;
    if (this.chatting !== previousValue) {
      if (this.historyView && this.historyView.list) {
        this.historyView.list = [];
      }
      if (!this.socket.histories) {
        return;
      }
      const data = this.socket.histories.get(this.socket.chattingUser);
      if (data) {
        this.chatList = data;
      } else {
        this.chatList = new Collections.LinkedList<ImHistoryMessage>();
      }
    }
  }

  send(ev: KeyboardEvent) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.returnValue = false;
    this.msg = this.msg.trim();
    if (this.msg.length >= 250) {
      console.log('max length 255');
      return;
    }
    if (this.msg.length > 0) {
      this.addChat();
    }
  }

  addChat() {
    const now = new Date().getTime();
    this.socket.msg(this.chatting, this.msg);
    this.chatList = this.socket.histories.get(this.socket.chattingUser);
    if (this.chatList === undefined) {
      this.chatList = new Collections.LinkedList<ImHistoryMessage>();
    }
    this.chatList.add({ From: this.socket.key, Msg: this.msg, Timestamp: now }, 0);
    this.socket.saveHistorys(this.chatting, this.chatList);
    this.msg = '';
  }
}
