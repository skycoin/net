import { Component, OnInit, ViewChild } from '@angular/core';
import { SocketService } from '../providers';
import { ImRecentItemComponent } from '../components';
@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  chatWindow = false;
  recent_list = [
    'Mary',
    'Lucien',
    'Steve',
    'LiLei',
    'Apple',
    'Box',
    'Test',
    'Green',
    'White']
  constructor(private socket: SocketService) {
  }
  ngOnInit(): void {
    if (this.socket.chattingUser) {
      this.chatWindow = true;
    }
    // this.socket.start();
  }

}
