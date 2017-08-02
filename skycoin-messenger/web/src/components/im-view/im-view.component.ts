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
import { SocketService, ImHistoryMessage } from '../../providers';
import { ToolService } from '../../providers/tool/tool.service';
import * as Collections from 'typescript-collections';
import { ImHistoryViewComponent } from '../im-history-view/im-history-view.component';
import * as emojione from 'emojione';


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
  constructor(public socket: SocketService, private tool: ToolService) { }

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

  send(msg: string) {
    msg = msg.trim();
    if (msg.length > 0) {
      this.addChat(msg);
    }
  }

  addChat(msg: string) {
    const now = new Date().getTime();
    this.socket.msg(this.chatting, msg);
    this.chatList = this.socket.histories.get(this.socket.chattingUser);
    if (this.chatList === undefined) {
      this.chatList = new Collections.LinkedList<ImHistoryMessage>();
    }
    this.chatList.add({ From: this.socket.key, Msg: msg, Timestamp: now }, 0);
    this.socket.saveHistorys(this.chatting, this.chatList);
    // tslint:disable-next-line:no-unused-expression
    this.socket.recent_list[this.socket.getRencentListIndex(this.chatting)].last = emojione.shortnameToImage(msg);
  }
}
