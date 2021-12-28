<?php

require_once __DIR__ . "/config.php";

use \GatewayWorker\Lib\Gateway;
use \Workerman\Lib\Timer;
use \Class\Tool;

class Events
{
    public static function onConnect($client_id)
    {
        $auth_timer_id = Timer::add(10, function ($client_id) {
            Gateway::closeClient($client_id, json_encode(['type' => 'auth_timeout']));
        }, array($client_id), false);
        Gateway::updateSession($client_id, array('auth_timer_id' => $auth_timer_id));
    }

    public static function onMessage($client_id, $message)
    {
        $recv = json_decode($message, true);
        $client_ip = $_SERVER['REMOTE_ADDR'] ?? '1.1.1.1';

        if (!$recv || !isset($recv['type'])) {
            Tool::out("$client_id 非法的请求,断开连接");
            Gateway::closeClient($client_id, json_encode(['type' => 'invalid_request']));
            return;
        }
        $type = $recv['type'];

        if (!in_array($type, ['ping', 'server_info', 'login'])) {
            Tool::out("$client_id $message"); //TODO debug
        }
        switch ($type) {
            case "ping":
                Gateway::sendToClient($client_id, json_encode(['type' => 'pong']));
                break;
            case "login":
                $token = $recv['token'] ?? '';
                $platform = $recv['platform'] ?? '';

                if ($platform == "server") {
                    //服务器节点
                    if ($token !== SERVERTOKEN) {
                        Tool::out("$client_id 非法的token,断开连接");
                        Gateway::closeClient($client_id, json_encode(['type' => 'invalid_token']));
                        return;
                    }
                    $tag = $recv['tag'] ?? 'unKnow Server';
                    $server_info = [
                        'platform' => 'server',
                        'tag' => $tag,
                        'ip' => $client_ip,
                        'online_time' => time(),
                    ];
                    Gateway::joinGroup($client_id, 'server');
                    Gateway::updateSession($client_id, $server_info); // 设置用户标签
                    Gateway::sendToClient($client_id, json_encode(['type' => 'login_success']));
                    Gateway::sendToGroup("web", json_encode([
                        'type' => 'server_online',
                        "data" => $server_info,
                    ]));
                    Tool::out("${client_id} ${tag} 连接到服务器");
                } else {
                    //网页
                    Gateway::joinGroup($client_id, 'web');
                    Gateway::updateSession($client_id, [
                        'platform' => 'web',
                        'ip' => $client_ip,
                    ]); // 设置连接者信息

                    // 发送当前连接者信息
                    $list = Gateway::getClientIdListByGroup('server');
                    $server_lists = [];
                    foreach ($list as $key => $server_client_id) {
                        $client_id_info = Gateway::getSession($server_client_id);
                        array_push($server_lists, [
                            "client_id" => $server_client_id,
                            "ip" => $client_id_info['ip'],
                            "tag" => $client_id_info['tag'],
                            "online_time" => $client_id_info['online_time']
                        ]);
                    }
                    Gateway::sendToCurrentClient(json_encode(['type' => 'login_success', 'server_list' => $server_lists]));
                    //Tool::out("${client_id} (web) 连接到服务器");
                }
                Timer::del(Gateway::getSession($client_id)['auth_timer_id']); //删除Timer
                break;
            case 'server_info': //服务器回传 基本信息
                $server_info = $recv['data'] ?? [];
                $server_info['tag'] = Gateway::getSession($client_id)['tag'];
                Gateway::sendToGroup("web", json_encode(['type' => 'server_info', 'client_id' => $client_id, 'server_info' => $server_info]));
                break;
            default:
                Tool::out("$client_id 未知的请求类型");
        }
    }

    public static function onClose($client_id)
    {
        $client_platform = $_SESSION['platform'] ?? null;
        if ($client_platform == 'server') {
            $tag = $_SESSION['tag'] ?? null;

            GateWay::sendToGroup("web", json_encode(['type' => 'server_offline', 'client_id' => $client_id]));
            Tool::out("${client_id} ${tag} 失去连接");
        }
    }
}


/*
    三个 Group

    server: 服务器节点
    web: 网页节点
*/