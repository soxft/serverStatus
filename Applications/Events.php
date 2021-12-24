<?php

require_once __DIR__ . "/config.php";

use \GatewayWorker\Lib\Gateway;
use \Workerman\Lib\Timer;
use \Class\Tool;

class Events
{
    public static function onConnect($client_id)
    {
        Tool::out("$client_id connected");
        $auth_timer_id = Timer::add(15, function ($client_id) {
            Gateway::closeClient($client_id, json_encode(['type' => 'auth_timeout']));
            Tool::out("$client_id > 认证超时,断开连接");
        }, array($client_id), false);
        Gateway::updateSession($client_id, array('auth_timer_id' => $auth_timer_id));
    }

    public static function onMessage($client_id, $message)
    {
        $msgData = json_decode($message, true);

        if (!$msgData || !isset($msgData['type'])) {
            Tool::out("$client_id > 非法的请求,断开连接");
            Gateway::closeClient($client_id, json_encode(['type' => 'invalid_request']));
            return;
        }
        $type = $msgData['type'];
        if ($type !== 'ping') {
            Tool::out("$client_id > $message");
        }
        switch ($type) {
            case "ping":
                Gateway::sendToClient($client_id, json_encode(['type' => 'pong']));
                break;
            case "login":
                $token = $msgData['token'] ?? '';
                $platform = $msgData['platform'] ?? '';
                if ($platform == "server") {
                    //登录类型为服务器节点
                    if ($token !== SERVERTOKEN) {
                        Tool::out("$client_id > 非法的token,断开连接");
                        Gateway::closeClient($client_id, json_encode(['type' => 'invalid_token']));
                        return;
                    }
                    Gateway::joinGroup($client_id, 'server');
                    Gateway::sendToClient($client_id, json_encode(['type' => 'login_success']));
                } else if ($platform == "web") {
                    //网页 状态展示
                    Gateway::joinGroup($client_id, 'web');
                    Gateway::sendToClient($client_id, json_encode(['type' => 'login_success']));
                } else {
                    Tool::out("$client_id > 非法的token,断开连接");
                    Gateway::closeClient($client_id, json_encode(['type' => 'invalid_token']));
                    return;
                }
                Timer::del(Gateway::getSession($client_id)['auth_timer_id']); //删除Timer
                break;
            default:
                Tool::out("$client_id > 未知的请求类型,断开连接");
                Gateway::closeClient($client_id, json_encode(['type' => 'unknow_request']));
        }
    }

    public static function onClose($client_id)
    {
        // GateWay::sendToAll("$client_id logout\r\n");
        Tool::out("$client_id > logout");
    }
}
