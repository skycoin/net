import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import 'rxjs/add/operator/map';

@Injectable()
export class ApiService {
  // private url = 'http://127.0.0.1:4998/';
  // private connUrl = this.url + '/conn/';
  constructor(private httpClient: HttpClient) { }
  getAllNode() {
    return this.httpClient.get('/conn/getAll').map((res: ConnsResponse) => {
      const connDatas: Array<ConnData> = [];
      res.conns.forEach((v, i) => {
        connDatas.push({ index: i + 1, key: v.key, type: v.type });
      });
      return connDatas;
    });
  }
}
export interface Conn {
  key?: string;
  type?: string;
}
export interface ConnsResponse {
  conns?: Array<Conn>;
}

export interface ConnData {
  index?: number;
  key?: string;
  type?: string;
}
