<?php

use \GatewayWorker\Lib\Gateway;

use \Class\Tool;

class Events
{
    /**
     * 当客户端连接时触发
     * 如果业务不需此回调可以删除onConnect
     * 
     * @param int $client_id 连接id
     */
    public static function onConnect($client_id)
    {
        Gateway::sendToClient($client_id, json_encode(['type' => 'request_passwd']));
        Tool::outPut("$client_id connected");
    }

    /**
     * 当客户端发来消息时触发
     * @param int $client_id 连接id
     * @param mixed $message 具体消息
     */
    public static function onMessage($client_id, $message)
    {
        // 向所有人发送 
        //Gateway::sendToAll("$message\r\n", exclude_client_id: [$client_id]);
        Tool::outPut("$client_id said: $message");
    }

    /**
     * 当用户断开连接时触发
     * @param int $client_id 连接id
     */
    public static function onClose($client_id)
    {
        // 向所有人发送 
       // GateWay::sendToAll("$client_id logout\r\n");
        Tool::outPut("$client_id logout");
    }
}
