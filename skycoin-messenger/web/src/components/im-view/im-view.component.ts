import {
  Component,
  OnInit,
  ViewEncapsulation,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
  AfterViewChecked,
  AfterViewInit
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
export class ImViewComponent implements OnInit, AfterViewChecked {
  chatList: Collections.LinkedList<ImHistoryMessage>;
  msg = '';
  @ViewChild(ImHistoryViewComponent) historyView: ImHistoryViewComponent;
  @Input() chatting = '';
  constructor(public socket: SocketService) { }

  ngOnInit() {
  }

  ngAfterViewChecked() {
    if (!this.socket.histories) {
      return;
    }
    const data = this.socket.histories.get(this.socket.chattingUser);
    if (data) {
      setTimeout(() => {
        this.historyView.list = data.toArray();
      }, 10)
    } else {
      this.chatList = new Collections.LinkedList<ImHistoryMessage>();
    }
  }
  send(ev: KeyboardEvent) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.returnValue = false;
    this.msg = this.msg.trim();
    if (this.msg.length > 0) {
      this.addChat();
    }
  }

  addChat() {
    this.socket.msg(this.chatting, this.msg);
    this.chatList = this.socket.histories.get(this.socket.chattingUser);
    if (this.chatList === undefined) {
      this.chatList = new Collections.LinkedList<ImHistoryMessage>();
    }
    this.chatList.add({ From: this.socket.key, Msg: this.msg }, 0);
    this.historyView.list = this.chatList.toArray();
    this.socket.saveHistorys(this.chatting, this.chatList);
    this.msg = '';
  }
}
