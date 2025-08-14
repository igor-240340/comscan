import { useState } from 'react';
import './App.css';
import { Greet } from "../wailsjs/go/main/App";

import { Layout, Table, Button, Space, Switch } from 'antd';

const { Sider, Content } = Layout;

function App() {
    const columns = [
        { title: 'name', dataIndex: 'name', key: 'name' },
        { title: 'usb', dataIndex: 'usb', key: 'usb' },
        { title: 'vid', dataIndex: 'vid', key: 'vid' },
        { title: 'pid', dataIndex: 'pid', key: 'pid' }
    ];
    const data = [
        { key: 1, name: 'COM1', usb: 'false', vid: "-", pid: "-" },
        { key: 2, name: 'COM2', usb: 'true', vid: '0x123', pid: '0x456' }
    ];

    return (
        <div id="App">
            <Layout style={{ height: '100vh' }}>
                {/* Left column */}
                <Sider width={400} style={{ background: '#fff', padding: '16px' }}>
                    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>

                        {/* Top controls */}
                        <div style={{ marginBottom: '16px' }}>
                            <Space>
                                Автоскан:
                                <Switch defaultChecked />
                                <Button type="primary">Обновить</Button>
                            </Space>
                        </div>

                        {/* Bottom table fills remaining space */}
                        <div style={{ flex: 1, overflow: 'auto' }}>
                            <Table
                                columns={columns}
                                dataSource={data}
                                pagination={false}
                                size="small"
                            />
                        </div>
                    </div>
                </Sider>

                {/* Right column */}
                <Content style={{ padding: '16px', background: '#fff' }}>
                    <Table
                        columns={[{ title: 'Log', dataIndex: 'log', key: 'log' }]}
                        dataSource={[
                            { key: 1, log: 'Application started' },
                            { key: 2, log: 'Connected to COM3' }
                        ]}
                        pagination={false}
                        size="small"
                    />
                </Content>
            </Layout>
        </div>
    )
}

export default App
