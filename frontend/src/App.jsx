import { useState, useEffect } from 'react';
import './App.css';

import { UpdatePortList } from '../wailsjs/go/main/App'

import { Layout, Table, Button, Space, Switch } from 'antd';

const { Sider, Content } = Layout;

function App() {
    const [portListData, setPortListData] = useState([]);
    const [autoScan, setAutoScan] = useState(false);

    let autoScanIntervalId;
    useEffect(() => {
        if (autoScan)
            autoScanIntervalId = setInterval(updatePortList, 1000);
        else {
            clearInterval(autoScanIntervalId);
            autoScanIntervalId = null;
        }

        return () => {
            clearInterval(autoScanIntervalId);
            autoScanIntervalId = null;
        };
    }, [autoScan]);

    async function updatePortList() {
        console.log("front: updatePortList()");

        let portList = await UpdatePortList();
        portList = portList.map((portInfo, index) => ({ ...portInfo, key: index }));
        setPortListData(portList);
        console.table(portList);
    }

    const columns = [
        { key: 0, title: 'Name', dataIndex: 'Name' },
        { key: 1, title: 'Usb', dataIndex: 'Usb' },
        { key: 2, title: 'Vid', dataIndex: 'Vid' },
        { key: 3, title: 'Pid', dataIndex: 'Pid' }
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
                                <Switch onChange={setAutoScan} />

                                <Button type="primary" onClick={updatePortList}>Обновить</Button>
                            </Space>
                        </div>

                        {/* Bottom table fills remaining space */}
                        <div style={{ flex: 1, overflow: 'auto' }}>
                            <Table
                                columns={columns}
                                dataSource={portListData}
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
