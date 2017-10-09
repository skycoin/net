import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

@Injectable()
export class ApiService {
  private url = 'http://127.0.0.1:5998/';
  private connUrl = this.url + 'conn/';
  constructor(private httpClient: HttpClient) {
    if (environment.production) {
      this.connUrl = '/conn/';
    }
  }
  getAllNode() {
    return this.httpClient.get(this.connUrl + 'getAll')
      .catch(err => Observable.throw(err));
  }
}
export interface Conn {
  key?: string;
  type?: string;
  send_bytes?: number;
  received_bytes?: number;
  last_ack_time?: number;
  start_time?: number;
}
export interface ConnsResponse {
  conns?: Array<Conn>;
}

export interface ConnData extends Conn {
  index?: number;
}
