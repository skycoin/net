import { Component, OnInit, ViewEncapsulation, Input, ViewChildren, QueryList, Output, EventEmitter } from '@angular/core';
import { ImRecentItemComponent } from '../im-recent-item/im-recent-item.component';
import { SocketService, UserService } from '../../providers';
import { MdDialog } from '@angular/material';
import { CreateChatDialogComponent } from '../create-chat-dialog/create-chat-dialog.component';

@Component({
  selector: 'app-im-recent-bar',
  templateUrl: './im-recent-bar.component.html',
  styleUrls: ['./im-recent-bar.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImRecentBarComponent implements OnInit {
  chatting = '';
  @ViewChildren(ImRecentItemComponent) items: QueryList<ImRecentItemComponent>;
  @Input() list = [];
  constructor(private socket: SocketService, private user: UserService, private dialog: MdDialog) { }
  ngOnInit() {
  }

  selectItem(item: ImRecentItemComponent) {
    // this.chatting.emit(item);
    item.info.unRead = 0;
    this.chatting = item.info.name;
    this.socket.chattingUser = item.info.name;
    const tmp = this.items.filter((el) => {
      return el.info.name !== item.info.name;
    });
    tmp.forEach(el => {
      el.active = false;
    });
  }

  openCreate(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const def = this.dialog.open(CreateChatDialogComponent, { position: { top: '10%' }, width: '350px' });
    def.afterClosed().subscribe(key => {
      if (key !== '' && key) {
        const icon = this.user.getRandomMatch();
        this.socket.recent_list.push({ name: key, last: '', icon: icon });
        this.socket.userInfo.set(key, { Icon: icon })
      }
    })
  }
}
