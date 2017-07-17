import { Component, OnInit, ViewChild } from '@angular/core';
import { SocketService, ImHistoryMessage } from '../providers';
import { ImRecentItemComponent } from '../components';
@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  chatWindow = false;
  constructor(private socket: SocketService) {
  }
  ngOnInit(): void {
    const testUser = window.location.search.slice(6);
    if (testUser !== '') {
      this.socket.recent_list.push(testUser);
    }
    if (this.socket.chattingUser) {
      this.chatWindow = true;
    }
    this.socket.start();
  }

}
