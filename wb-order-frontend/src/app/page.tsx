"use client"
import { useState } from 'react'
import axios from 'axios'

export default function Home() {
  const [orderUid, setOrderUid] = useState('')
  const [orderData, setOrderData] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const fetchOrder = async () => {
    if (!orderUid.trim()) {
      setError('Введите order_uid')
      return
    }

    setLoading(true)
    setError('')

    try {
      const response = await axios.get(`http://localhost:8081/order/${orderUid}`)
      setOrderData(response.data.order)
    } catch (err) {
      setError('Заказ не найден')
      setOrderData(null)
    } finally {
      setLoading(false)
    }
  }

  return (
      <div className="min-h-screen bg-gray-100 py-[25%] px-4">
        <div className="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden p-6">
          <h1 className="text-2xl font-bold text-gray-800 mb-6">Поиск заказа</h1>

          <div className="flex mb-6">
            <input
                type="text"
                value={orderUid}
                onChange={(e) => setOrderUid(e.target.value)}
                placeholder="Введите order_uid"
                className="flex-grow px-4 py-2 border text-black rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
                onClick={fetchOrder}
                disabled={loading}
                className="bg-purple-600 text-white px-4 py-2 rounded-r-lg hover:purple disabled:bg-purple-400"
            >
              {loading ? 'Загрузка...' : 'Найти'}
            </button>
          </div>

          {error && <p className="text-red-500 mb-4">{error}</p>}

          {orderData && (
              <div className="border-t pt-4">
                <h2 className="text-xl text-black font-semibold mb-3">Данные заказа:</h2>
                <pre className="bg-gray-100 text-black p-4 rounded-lg overflow-auto text-sm">
              {JSON.stringify(orderData, null, 2)}
            </pre>
              </div>
          )}
        </div>
      </div>
  )
}