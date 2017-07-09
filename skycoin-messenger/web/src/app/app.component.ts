import { Component,OnInit } from '@angular/core';
import { SocketService } from '../providers';

@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  
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
    this.socket.start();
    setTimeout(() => {
      console.log('send msg');
      this.socket.msg('hi !!!!');
    },3000)
  }
}
